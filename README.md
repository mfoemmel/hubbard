Description
===========

Hubbard is a command line tool for managing shared git repositories in a team environment. It provides a security model so that users can easily control who has access to the projects they create. 

Hubbard uses public SSH keys to keep track of who is executing what commands. This means you only have to create a single account on the server, instead of one per user.

Hubbard was heavily inspired by gitosis, another tool for managing git repositories. However, the goal of Hubbard was to place less burden on the system administrator by allowing users to manage permissions for their own projects.

How It Works
============

All comminication between users and the Hubbard server happens over SSH. Users must register their public SSH keys with the server before they can connect to it.

When a user connects to the Hubbard server, the SSH daemon tries to find the user's public SSH key the "~/.ssh/authorized_keys" file on the server. That file also contains information about which user to associate with that SSH key. That information is automatically passed to the "hubbard" executable, so there is no way for users to run other programs on the server. 

Installation
===========

### Server ###

The first step is to create a user account called "hub" on the server machine and log into that account. You'll also need to make sure that Ruby and Rubygems are installed on the machine. You can then install hubbard like this:

    $ gem install hubbard

Unless you installed the gem using "sudo", the "hubbard" executable will be found somewhere under the "~/.gem" directory. You'll need to make sure this directory is included in the PATH whenever anyone uses SSH to connect. You can do this by adding the following line to the top of your "~/.bashrc" file (make sure the path matches):

    export PATH=$PATH:~/.gem/ruby/1.8/bin

Create the directories and files that SSH needs to work.
    [in hub's home directory]
    mkdir .ssh
    chmod 700 .ssh
    touch .ssh/authorized_keys
    chmod 600 .ssh/authorized_keys

The next step is to create an SSH keypair to access the "admin" account on the Hubbard server. You should only use this key when performing tasks that require admin access. Run this on the machine that you'll be accessing Hubbard from (i.e. your local workstation, not the server):

    $ ssh-keygen -f ~/.ssh/hubadmin

Now we'll need to copy the public key up to the server. Assuming that the hubbard server is named "example" :

    $ cat ~/.ssh/hubadmin.pub | ssh hub@example hubbard admin add-key default

On your local workstation, you can now create an alias for executing admin commands:

    $ alias hubadmin='ssh -i ~/.ssh/hubadmin hub@example.com'

Test the configuration:

    $ hubadmin list-users
    admin

You'll probably want to create a normal (i.e. non-admin) account as described in the next section.

### Client ###

There is no executable to install on the client, but you will have to generate an SSH keypair and send the public key (i.e. the one that ends with ".pub") to your administrator. Your administrator will have to run the following command: 

    $ cat <keyfile> | hubadmin run-as <username> add-key default

Assuming your SSH keys have been set up correcly, you can simply SSH into the server machine to executable commands. You'll probably want to set up an alias to make this easier, however. For example, if your server is running at "example", you can do this:

    $ alias hub='ssh hub@example'

To test it, run:
    
    $ hub help

