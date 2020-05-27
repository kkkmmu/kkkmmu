package main

import (
	"flag"
	"fmt"
	"github.com/marcusolsson/tui-go"
	"github.com/sahilm/fuzzy"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"sync"
)

var IP = flag.String("ip", "", "IP address of the remote device")
var (
	method  = "GET"
	params  = "x=2&y=3"
	payload = `{"id": 12}`
	headers = "User-Agent: myBrowser"

	FnNameMap          = make(map[string]*Fn, 1000)
	FnAddrMap          = make(map[string]*Fn, 1000)
	Fns                = make([]string, 0, 1000)
	Ads                = make([]string, 0, 1000)
	MFns               []*Fn
	MFnsLock           sync.Mutex
	FNR                = regexp.MustCompile("(?P<addr>[[:alnum:]]+) <(?P<name>[[:alnum:]_]+)>:")
	CalleeR            = regexp.MustCompile("<(?P<name>[[:alnum:]_@]+)>")
	NoMatchStr         = "No matching for %s!"
	CHistory           = make([]string, 0, 100)
	NameToToolChainMap = map[string]string{
		"M3000":  "/opt/freescale/usr/local/gcc-4.3.74-eglibc-2.8.74-dp-2/powerpc-none-linux-gnuspe/bin/powerpc-none-linux-gnuspe-",
		"V3306G": "/opt/eldk-5.5/armv7a-softfloat-p1/sysroots/i686-linux/usr/bin/armv7a-linux-gnueabi/arm-linux-gnueabi-",
		"M2400":  "/opt/freescale/usr/local/gcc-4.3.74-eglibc-2.8.74-dp-2/powerpc-none-linux-gnuspe/bin/powerpc-none-linux-gnuspe-",
		"M2300":  "/opt/freescale/usr/local/gcc-4.3.74-eglibc-2.8.74-dp-2/powerpc-none-linux-gnuspe/bin/powerpc-none-linux-gnuspe-",
	}
	Product   = flag.String("product", "M3000", "Product that the bin file compiled for")
	ToolChain = flag.String("toolchain", "", "ToolChain which can be used to analysis the bin file")
	Bin       = flag.String("bin", "", "The bin file name to analysis")
)

type Fn struct {
	Name   string
	Addr   string
	Code   string
	Callee []string
}

func main() {
	asmLbl := tui.NewLabel("")
	asmLbl.SetSizePolicy(tui.Expanding, tui.Expanding)

	scroll := tui.NewScrollArea(asmLbl)

	respBody := tui.NewVBox(scroll)
	respBody.SetTitle("Function")
	respBody.SetBorder(true)

	clist := tui.NewList()

	clist.OnSelectionChanged(func(l *tui.List) {
		if l.Selected() >= 0 && l.Selected() <= len(FnNameMap) {
			if f, ok := FnNameMap[clist.SelectedItem()]; ok {
				asmLbl.SetText(f.Code)
			} else {
				asmLbl.SetText(fmt.Sprintf("No assembly code for %s", clist.SelectedItem()))
			}
			scroll.ScrollToTop()
		} else {
			asmLbl.SetText(fmt.Sprintf("Invalid index: %d", l.Selected))
		}
	})

	clist.OnItemActivated(func(l *tui.List) {
		if l.Selected() >= 0 && l.Selected() <= len(MFns) {
			f, ok := FnNameMap[clist.SelectedItem()]
			if ok {
				asmLbl.SetText(f.Code)
			} else {
				asmLbl.SetText(fmt.Sprintf("No assembly code for %s", clist.SelectedItem()))
				return
			}
			scroll.ScrollToTop()
			clist.RemoveItems()
			clist.AddItems(f.Callee...)
			clist.Select(0)
			CHistory = append(CHistory, clist.SelectedItem())

		} else {
			asmLbl.SetText(fmt.Sprintf("Invalid index: %d", l.Selected))
		}
	})

	cbox := tui.NewHBox(clist)
	cbox.SetTitle("Callee")
	cbox.SetBorder(true)

	mlist := tui.NewList()

	//mlist.SetTitile("Matches")
	//mlist.SetBorder(true)
	mlist.OnSelectionChanged(func(l *tui.List) {
		if len(CHistory) > 0 {
			CHistory = make([]string, 0, 1000)
		}

		if l.Selected() >= 0 && l.Selected() <= len(MFns) {
			MFnsLock.Lock()
			asmLbl.SetText(MFns[l.Selected()].Code)
			MFnsLock.Unlock()
			scroll.ScrollToTop()

			clist.RemoveItems()
			clist.AddItems(MFns[l.Selected()].Callee...)
			clist.Select(0)
			CHistory = append(CHistory, clist.SelectedItem())
		} else {
			asmLbl.SetText(fmt.Sprintf("Invalid index: %d", l.Selected))
		}

	})

	mbox := tui.NewHBox(mlist)
	mbox.SetTitle("Matches")
	mbox.SetBorder(true)
	mcbox := tui.NewHBox(mbox, cbox)
	match := tui.NewHBox(mcbox, respBody)

	match.SetSizePolicy(tui.Expanding, tui.Preferred)

	browser := tui.NewVBox(match)
	browser.SetSizePolicy(tui.Preferred, tui.Expanding)

	input := tui.NewEntry()
	input.SetText("main")
	input.SetSizePolicy(tui.Preferred, tui.Expanding)
	input.OnChanged(func(e *tui.Entry) {
		if e.Text() == "" {
			return
		}

		if strings.HasPrefix(e.Text(), "x/ ") {
			addr := strings.TrimSpace(strings.TrimPrefix(e.Text(), "x/ "))
			matches := fuzzy.Find(addr, Ads)
			if len(matches) == 0 {
				//respHeadLbl.SetText(fmt.Sprintf(NoMatchStr, e.Text()))
				asmLbl.SetText(fmt.Sprintf(NoMatchStr, e.Text()))
				return
			}

			ms := make([]string, 0, len(matches))

			MFnsLock.Lock()
			MFns = make([]*Fn, 0, len(matches))
			for _, m := range matches {
				ms = append(ms, m.Str)
				MFns = append(MFns, FnAddrMap[m.Str])
			}
			MFnsLock.Unlock()
			//respHeadLbl.SetText(strings.Join(ms, "\n"))
			if _, ok := FnAddrMap[ms[0]]; ok {
				asmLbl.SetText(FnAddrMap[ms[0]].Code)
			} else {
				asmLbl.SetText(fmt.Sprintf(NoMatchStr, e.Text()))
			}

			mlist.RemoveItems()

			if len(ms) > 20 {
				mlist.AddItems(ms[0:20]...)
			} else {
				mlist.AddItems(ms...)
			}
			mlist.Select(0)

		} else {
			matches := fuzzy.Find(e.Text(), Fns)
			if len(matches) == 0 {
				//respHeadLbl.SetText(fmt.Sprintf(NoMatchStr, e.Text()))
				asmLbl.SetText(fmt.Sprintf(NoMatchStr, e.Text()))
				return
			}
			ms := make([]string, 0, len(matches))

			MFnsLock.Lock()
			MFns = make([]*Fn, 0, len(matches))
			for _, m := range matches {
				ms = append(ms, m.Str)
				MFns = append(MFns, FnNameMap[m.Str])
			}
			MFnsLock.Unlock()
			//respHeadLbl.SetText(strings.Join(ms, "\n"))
			if _, ok := FnNameMap[ms[0]]; ok {
				asmLbl.SetText(FnNameMap[ms[0]].Code)
			} else {
				asmLbl.SetText(fmt.Sprintf(NoMatchStr, e.Text()))
			}

			mlist.RemoveItems()

			if len(ms) > 20 {
				mlist.AddItems(ms[0:20]...)
			} else {
				mlist.AddItems(ms...)
			}
			mlist.Select(0)
		}
	})

	inBox := tui.NewHBox(input)
	inBox.SetTitle("Function/Address")
	inBox.SetBorder(true)

	root := tui.NewVBox(inBox, browser)

	tui.DefaultFocusChain.Set(input, mlist, clist, scroll)

	theme := tui.NewTheme()
	theme.SetStyle("box.focused.border", tui.Style{Fg: tui.ColorYellow, Bg: tui.ColorDefault})
	theme.SetStyle("list.item.selected", tui.Style{Fg: tui.ColorGreen, Bg: tui.ColorDefault, Underline: tui.DecorationOn})

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	ui.SetTheme(theme)
	ui.SetKeybinding("Esc", func() { ui.Quit() })
	ui.SetKeybinding("Up", func() {
		if scroll.IsFocused() {
			scroll.Scroll(0, -1)
		}
	})
	ui.SetKeybinding("Down", func() {
		if scroll.IsFocused() {
			scroll.Scroll(0, 1)
		} /*else if mlist.IsFocused() {
			if mlist.Selected() < mlist.Length() {
				mlist.OnKeyEvent(tui.KeyEvent{Key: tui.KeyDown})
			}
		}*/
	})
	ui.SetKeybinding("Left", func() {
		if scroll.IsFocused() {
			scroll.Scroll(-1, 0)
		} else if clist.IsFocused() {
			if len(CHistory) > 1 {
				c := CHistory[len(CHistory)-1]
				CHistory = CHistory[:len(CHistory)-1]
				f, ok := FnNameMap[c]
				if ok {
					asmLbl.SetText(f.Code)
					scroll.ScrollToTop()
					clist.RemoveItems()
					clist.AddItems(f.Callee...)
					clist.Select(0)
				}
			} else {
				f, ok := FnNameMap[CHistory[0]]
				if ok {
					asmLbl.SetText(f.Code)
					scroll.ScrollToTop()
					clist.RemoveItems()
					clist.AddItems(f.Callee...)
					clist.Select(0)
				}
			}
		}

	})
	ui.SetKeybinding("Right", func() {
		if scroll.IsFocused() {
			scroll.Scroll(1, 0)
		}
	})
	ui.SetKeybinding("a", func() {
		if scroll.IsFocused() {
			scroll.SetAutoscrollToBottom(true)
		}
	})
	ui.SetKeybinding("t", func() {
		if scroll.IsFocused() {
			scroll.ScrollToTop()
		}
	})
	ui.SetKeybinding("b", func() {
		if scroll.IsFocused() {
			scroll.ScrollToBottom()
		}
	})

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	flag.Parse()
	var tool string
	if *ToolChain != "" {
		tool = *ToolChain
	} else {
		if t, ok := NameToToolChainMap[*Product]; ok {
			tool = t
		} else {
			panic("You must give the product name or tool chain")
		}
	}

	if *Bin == "" {
		panic("You must give the bin file to analsysi")
	}

	cmd := exec.Command(tool+"objdump", "-d", *Bin)
	res, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	fs := strings.Split(string(res), "\n\n")

	for _, f := range fs {
		if strings.TrimSpace(f) == "" {
			continue
		}
		mts := FNR.FindStringSubmatch(f)
		if len(mts) != 3 {
			//fmt.Println(f)
			continue
		}
		fn := &Fn{
			Name: mts[2],
			Addr: mts[1],
			Code: f,
		}

		FnNameMap[mts[2]] = fn
		FnAddrMap[mts[1]] = fn
		Fns = append(Fns, mts[2])
		Ads = append(Ads, mts[1])
	}

	for _, fn := range FnNameMap {
		cs := CalleeR.FindAllStringSubmatch(fn.Code, -1)
		fn.Callee = make([]string, 0, len(cs))
		for _, c := range cs {
			fn.Callee = append(fn.Callee, c[1])
		}
	}
}
