package topology

import (
	"dut"
	"fmt"
	"ifp"
	"its"
	"link"
	"peer"
)

type Topology struct {
	DUTs  map[string]*dut.DUT
	Links map[*peer.Peer]*link.Link
	ITSs  map[string]*its.ITS
}

func New() *Topology {
	return &Topology{
		DUTs:  make(map[string]*dut.DUT, 3),
		ITSs:  make(map[string]*its.ITS, 3),
		Links: make(map[*peer.Peer]*link.Link, 10),
	}
}

func (t *Topology) Init() error {
	for _, d := range t.DUTs {
		err := d.Connect()
		if err != nil {
			return fmt.Errorf("Topology Init failed, cannot connect to %s with %s", d.Name, err)
		}
		fmt.Printf("Connect to %s success!", d.Name)
	}

	for _, i := range t.ITSs {
		err := i.Connect()
		if err != nil {
			return fmt.Errorf("Topology Init failed, cannot connect to %s with %s", i.Name, err)
		}
		fmt.Printf("Connect to %s success!", i.Name)
	}

	return nil
}

func (t *Topology) Verify() error {
	for _, d := range t.DUTs {
		if err := d.Verify(); err != nil {
			return fmt.Errorf("Invalid topology with %s", err)
		}
	}

	for _, i := range t.ITSs {
		if err := i.Verify(); err != nil {
			return fmt.Errorf("Invalid topology with %s", err)
		}
	}

	for _, l := range t.Links {
		itf1 := t.GetInteraceByName(l.P1.Name)
		if itf1 == nil {
			return fmt.Errorf("Cannot find interface %s in this topology", l.P1.Name)
		}

		itf2 := t.GetInteraceByName(l.P2.Name)
		if itf2 == nil {
			return fmt.Errorf("Cannot find interface %s in this topology", l.P2.Name)
		}

		itf1.LinkTo = itf2
		itf2.LinkTo = itf1
	}

	for _, d := range t.DUTs {
		for _, itf := range d.Ifps {
			if itf.LinkTo == nil {
				return fmt.Errorf("There is no link on interface %s of %s", itf.Name, d.Name)
			}
		}
	}

	return nil
}

func (t *Topology) GetInteraceByName(name string) *ifp.Ifp {
	for _, d := range t.DUTs {
		if itf, ok := d.IfpMap[name]; ok {
			return itf
		}
	}

	for _, i := range t.ITSs {
		if itf, ok := i.IfpMap[name]; ok {
			return itf
		}
	}

	return nil
}

func (t *Topology) AddDUT(nd *dut.DUT) error {
	if _, ok := t.DUTs[nd.Name]; ok {
		return fmt.Errorf("DUT %s alread exist", nd.Name)
	}

	t.DUTs[nd.Name] = nd

	return nil
}

func (t *Topology) DelDUT(dd *dut.DUT) error {
	if _, ok := t.DUTs[dd.Name]; !ok {
		return fmt.Errorf("DUT %s does not exist", dd.Name)
	}

	delete(t.DUTs, dd.Name)

	return nil
}

func (t *Topology) AddTester(nits *its.ITS) error {
	if _, ok := t.ITSs[nits.Name]; ok {
		return fmt.Errorf("Tester %s alread exist", nits.Name)
	}

	t.ITSs[nits.Name] = nits

	return nil
}

func (t *Topology) DelTester(dits *its.ITS) error {
	if _, ok := t.ITSs[dits.Name]; !ok {
		return fmt.Errorf("Tester %s does not exist", dits.Name)
	}

	delete(t.ITSs, dits.Name)

	return nil
}

func (t *Topology) AddLink(newlink *link.Link) error {
	t.Links[newlink.P1] = newlink
	t.Links[newlink.P2] = newlink

	return nil
}

func (t *Topology) DelLink(oldlink *link.Link) error {
	delete(t.Links, oldlink.P1)
	delete(t.Links, oldlink.P2)

	return nil
}

func (t *Topology) GetInterfaceCount() int {
	var count int
	for _, d := range t.DUTs {
		count += len(d.Ifps)
	}

	for _, i := range t.ITSs {
		count += len(i.Ifps)
	}

	return count
}
