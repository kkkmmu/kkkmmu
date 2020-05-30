Vi mode in Bash
Posted on 2012-02-28
Because the Bash shell uses GNU Readline, you’re able to set options in .bashrc and .inputrc to extensively customise how you enter your command lines, including specifying whether you want the default Emacs-like keybindings (e.g. Ctrl+W to erase a word), or bindings using a modal interface familiar to vi users.

To use vi mode in Bash and any other tool that uses GNU Readline, such as the MySQL command line, you need only put this into your .inputrc file, then log out and in again:

set editing-mode vi
If you only want to use this mode in Bash, an alternative is to use the following in your .bashrc, again with a login/logout or resourcing of the file:

set -o vi
You can check this has applied correctly with the following, which will spit out a list of the currently available bindings for a set of fixed possible actions for text:

$ bind -P
The other possible value for editing-mode is emacs. Both options simply change settings in the output of bind -P above.

This done, you should find that you open a command line in insert mode, but when you press Escape or Ctrl+[, you will enter an emulation of Vi’s normal mode. Enter will run the command in its present state while in either mode. The bindings of most use in this mode are the classic keys for moving to positions on a line:

^ — Move to start of line
$ — Move to end of line
b — Move back a word
w — Move forward a word
e — Move to the end of the next word
Deleting, yanking, and pasting all work too. Also note that k and j in normal mode allow you to iterate through your history in the same way that Ctrl+P and Ctrl+N do in Emacs mode. Searching with ? and / doesn’t work by default, but the Ctrl+R and Ctrl+S bindings still seem to work.

If you are reasonably used to working with the modeless Emacs for most of your shell work but do occasionally find yourself in a situation where being able to edit a long command line in vi would be handy, you may prefer to press Ctrl+X Ctrl+E to bring up the command line in your $EDITOR instead, affording you the complete power of Vim to edit rather than the somewhat sparse vi emulation provided by Readline. If you want to edit a command you just submitted, the fc (fix command) Bash builtin works too.

The shell is a very different environment for editing text, being inherently interactive and fixed to one line rather than a static multi-line buffer, which leads some otherwise diehard vi users such as myself to prefer the Emacs mode for editing commands. For one thing, vi mode in Bash trips on the vi anti-pattern of putting you in insert mode by default, and at the start of every command line. But if the basic vi commands are etched indelibly into your memory and have become automatic, you may appreciate being able to edit your command line using these keys instead. It’s certainly worth a try at any rate, and I do still occasionally see set -o vi in the .bashrc files of a few of the pros on GitHub.
