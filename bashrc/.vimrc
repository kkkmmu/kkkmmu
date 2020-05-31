"Install vim-plug
if empty(glob('~/.vim/autoload/plug.vim'))
  silent !curl -fLo ~/.vim/autoload/plug.vim --create-dirs
    \ https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim
  autocmd VimEnter * PlugInstall --sync | source $MYVIMRC
endif

"vim plugin packages
"wget https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim
"PlugInstall: Install plugins
"PlugUpdate:  Install or update plugins
"PlugClean:   Remove unlisted plugins
"PlugUpgrade: Ugrade vim-plug itself
"PlugStatus:  Check the status of plugins
"rtp:         Subdirectory that contains VIM plugin
"dir:         Custom directory for the plugin
"as:          Use different name for the plugin
"do:          Post-update hook
"on:          On-demand loading: Comannds for <Plug>-mappings
"for:         On-demand mapping: File types
"frozen:      Do not update unless explictly specified
call plug#begin('~/.vim/plugged')
Plug 'tpope/vim-sensible'
Plug 'junegunn/seoul256.vim'
Plug 'Shougo/vimproc.vim', { 'do': 'make' }
Plug 'scrooloose/nerdtree', { 'on' : 'NERDTreeToggle' }
Plug 'junegunn/goyo.vim', { 'for' : 'markdown' }
Plug 'junegunn/limelight.vim'
Plug 'rdnetto/YCM-Generator', { 'branch': 'stable' }
Plug 'ycm-core/YouCompleteMe', { 'do': './install.py --all' }
Plug 'nsf/gocode', { 'tag': '*', 'rtp': 'vim' }
Plug 'fatih/vim-go', { 'do': ':GoInstallBinaries' }
Plug 'tpope/vim-fugitive'
Plug 'vim-scripts/grep.vim'
Plug 'tell-k/vim-autopep8', {'for' : 'py'}
Plug 'nvie/vim-flake8'
Plug 'honza/vim-snippets'
Plug 'vim-scripts/indentpython.vim',  {'for' : 'py'}
Plug 'Lokaltog/powerline', {'rtp': 'powerline/bindings/vim/'}
Plug 'nathanalderson/yang.vim', { 'for' : 'yang' }
Plug 'vim-utils/vim-man'
Plug 'nrc/rustfmt', {'for' : 'rs'}
Plug 'rust-lang/rust.vim', {'for' : 'rs'}
Plug 'pangloss/vim-javascript', {'for' : 'js'}
"Plug 'w0rp/ale'
Plug 'flazz/vim-colorschemes'
Plug 'ludovicchabant/vim-gutentags'
Plug 'skywind3000/gutentags_plus'
Plug 'neomake/neomake'
"Plug 'kien/ctrlp.vim'
Plug 'Yggdroot/LeaderF', { 'do': './install.sh' }
Plug 'Shougo/vimshell.vim'
Plug 'Shougo/vimproc.vim'
Plug 'vim-syntastic/syntastic'
if has('nvim')
	Plug 'Shougo/deoplete.nvim', { 'do': ':UpdateRemotePlugins' }
else
	Plug 'Shougo/deoplete.nvim'
	Plug 'roxma/nvim-yarp'
	Plug 'roxma/vim-hug-neovim-rpc'
endif
let g:deoplete#enable_at_startup = 1
call plug#end() 

set nu
set hlsearch
set incsearch  "start search before pressing entry
"set listchars  "Makes set list prettier
"set scrolloff  "Aways show at least one line above/below the cursor
set autoread
set ruler
set cindent
set laststatus=2
set nolist
set noexpandtab
set tabstop=4
set linespace=4
set shiftwidth=4
set updatetime=250
set hidden
set t_Co=256
set autoindent
set smartindent
set smarttab
set foldmethod=manual
set encoding=utf-8

filetype plugin indent on
syntax on

let mapleader="."

map <leader>k <Plug>(Man)
map <leader>v <Plug>(Vman)
map <C-q> :NERDTreeToggle<CR>

nnoremap <leader>vs :VimShell<CR>
nnoremap <leader>vc :VimShellClose<CR>
nnoremap <leader>vp :VimShellPop<CR>
nnoremap <leader>d :GoDoc<CR>

"nnoremap <space> za
"
"egrep key mapping currently not used
"nnoremap <leader>rg :Rgrep<CR>

"LeaderF rg need install ripgrep
nnoremap <leader>rc :Leaderf --recall<CR>
nnoremap <leader>rg :<C-U><C-R>=printf("Leaderf! --stayOpen rg -e %s", expand("<cword>"))<CR>

"For git
nnoremap <leader>gs :Git<CR>
nnoremap <leader>gl :Git log<CR>
nnoremap <leader>gd :Git diff<CR>
nnoremap <leader>gb :Git blame<CR>

let Grep_Default_Filelist='*.[cshS] *.go *.py *.cpp *.hpp *.cc *.js'
let Grep_Skip_Files = '*.bak *~ *.so *.i *.a *.o'
let Grep_Default_Options='-i --color=auto'
let Grep_Skip_Dirs='.svn .git'

"let g:Lf_WindowPosition = 'popup'
let g:Lf_PreviewInPopup = 1
let g:Lf_ShortcutF = '<C-P>'
let g:Lf_ReverseOrder = 1
let g:Lf_AutoResize = 1
"let g:Lf_WorkingDirectoryMode = 'Ac'
"let g:Lf_RootMarkers = ['.git', '.svn']
"let g:Lf_WorkingDirectory = finddir('.git', '.svn', '.;')
let g:Lf_PreviewCode = 1

"pip install pygments is necessary
"let g:Lf_Gtagslabel = 'native-pygments'
let g:Lf_CommandMap = {'<C-K>': ['<Up>'], '<C-J>': ['<Down>']}
let g:Lf_UseVersionControlTool = 0
"let g:Lf_RgConfig = [
"        \ "--max-columns=150",
"        \ "--type-add web:*.{html,css,js}*",
"        \ "--glob=!git/*",
"        \ "--hidden"
"    \ ]

let g:autopep8_on_save = 1
filetype indent on
 let g:autopep8_disable_show_diff=1

"gutentags configuration
let g:gutentags_project_root = ['package.json', '.git', '.svn']
let g:gutentags_add_default_project_roots = 0
let g:gutentags_cache_dir = expand('~/.cache/vim/ctags/')
let g:gutentags_generate_on_write = 1
let g:gutentags_generate_on_new = 1
let g:gutentags_generate_on_missing = 1
let g:gutentags_generate_on_write = 1
let g:gutentags_generate_on_empty_buffer = 0
let g:gutentags_ctags_extra_args = [
      \ '--tag-relative=yes',
      \ '--fields=+ailmnS',
      \ ]
 
"gutentags ignored file
let g:gutentags_ctags_exclude = [
      \ '*.git', '*.svg', '*.hg',
      \ '*/tests/*',
      \ 'build',
      \ 'dist',
      \ '*sites/*/files/*',
      \ 'bin',
      \ 'node_modules',
      \ 'bower_components',
      \ 'cache',
      \ 'compiled',
      \ 'docs',
      \ 'example',
      \ 'bundle',
      \ 'vendor',
      \ '*.md',
      \ '*-lock.json',
      \ '*.lock',
      \ '*bundle*.js',
      \ '*build*.js',
      \ '.*rc*',
      \ '*.json',
      \ '*.min.*',
      \ '*.map',
      \ '*.bak',
      \ '*.zip',
      \ '*.pyc',
      \ '*.class',
      \ '*.sln',
      \ '*.Master',
      \ '*.csproj',
      \ '*.tmp',
      \ '*.csproj.user',
      \ '*.cache',
      \ '*.pdb',
      \ 'tags*',
      \ 'cscope.*',
      \ '*.css',
      \ '*.less',
      \ '*.scss',
      \ '*.exe', '*.dll',
      \ '*.mp3', '*.ogg', '*.flac',
      \ '*.swp', '*.swo',
      \ '*.bmp', '*.gif', '*.ico', '*.jpg', '*.png',
      \ '*.rar', '*.zip', '*.tar', '*.tar.gz', '*.tar.xz', '*.tar.bz2',
      \ '*.pdf', '*.doc', '*.docx', '*.ppt', '*.pptx',
      \ ]

map <leader>f  :YcmCompleter GoToDefinitionElseDeclaration<CR>
let g:ycm_autoclose_preview_window_after_completion=1

" Set this variable to 1 to fix files when you save them.
"let g:ale_fix_on_save = 1
"let g:ale_completion_enabled = 1
"let g:ale_completion_tsserver_autoimport = 1
"let g:ale_set_balloons = 1
"let g:ale_hover_to_preview = 1

set statusline+=%#warningmsg#
set statusline+=%{SyntasticStatuslineFlag()}
set statusline+=%*

let g:syntastic_always_populate_loc_list = 1
let g:syntastic_auto_loc_list = 1
let g:syntastic_check_on_open = 1
let g:syntastic_check_on_wq = 0

" enable gtags module
let g:gutentags_modules = ['ctags', 'gtags_cscope']

" config project root markers.
let g:gutentags_project_root = ['.git', '.root', '.svn']

" generate datebases in my cache directory, prevent gtags files polluting my project
let g:gutentags_cache_dir = expand('~/.cache/tags')
"
" change focus to quickfix window after search (optional).
let g:gutentags_plus_switch = 1

"let g:gutentags_define_advanced_commands = 1
"noremap <silent> <leader>u :GutentagsUpdate<C-R>

".noremap <silent> <leader>gs :GscopeFind s <C-R><C-W><cr>
".noremap <silent> <leader>gg :GscopeFind g <C-R><C-W><cr>
".noremap <silent> <leader>gc :GscopeFind c <C-R><C-W><cr>
".noremap <silent> <leader>gt :GscopeFind t <C-R><C-W><cr>
".noremap <silent> <leader>ge :GscopeFind e <C-R><C-W><cr>
".noremap <silent> <leader>gf :GscopeFind f <C-R>=expand("<cfile>")<cr><cr>
".noremap <silent> <leader>gi :GscopeFind i <C-R>=expand("<cfile>")<cr><cr>
".noremap <silent> <leader>gd :GscopeFind d <C-R><C-W><cr>
".noremap <silent> <leader>ga :GscopeFind a <C-R><C-W><cr>
".noremap <silent> <leader>gz :GscopeFind z <C-R><C-W><cr>

syntax on

if has('gui_running')
	set background=light
    "colorscheme solarized
else
	set background=dark
    "colorscheme solarized
	"colorscheme zenburn
endif

"Mark extra withespace as bad and probably color it red.
au BufRead,BufNewFile *.py,*.pyw,*.c,*.h match BadWhitespace /\s\+$/
highlight BadWhitespace ctermbg=red guibg=darkred

colorscheme molokai
"colorscheme koehler
"colorscheme evening
"colorscheme morning
"colorscheme pablo
"colorscheme ron
"colorscheme shine
"colorscheme torte
"colorscheme zeliner
"colorscheme hybrid
"colorscheme jellybeans
"colorscheme PaperColor
"colorscheme desert
"colorscheme desert256
"colorscheme airline

"autocmd fileType c,cpp ClangFormatAutoEnable
"autocmd vimenter *  TlistOpen

"dw, -----> delete word
"yt, -----> copy until ',' ---- we can do yt anything.
"yw, -----> copy word.
"dt, -----> delete until ',' ---- we can do dt anything.
"ge, -----> end.
"gE, -----> Delimit end.
"Tx, -----> Till x
"tx, -----> till x
"b, -----> word
"B, -----> word
"Fx, -----> Find x ',' ---> previous x ';' ----> next x
"fx, -----> find x ',' ---> previous x ';' ----> next x
"L, -----> Screen
"H, -----> Screen
"#, -----> find word under cousor
"*, -----> find word under cousor
"'', -----> last location
"'., -----> last edit
"dib -----> delete all content in ()
"yap -----> copy current paragrap
"dap -----> Delete all current paragraph
"<leader>di ----------> start DrawIt
"<leader>ds ----------> stop DrawIt

