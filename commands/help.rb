$stderr.puts <<-END
Usage: hub <command>

Projects:

list-projects
create-project <project>
delete-project <project>

Repositories:

list-repositories <project>
create-repository <project> <repository>
delete-repository <project> <repository>
fork-repository <from-project> <from-repository> <to-project> <to-repository>

Permissions:

list-permissions <project>
add-permission <project> <username> read|write|admin
remove-permission <project> <username>

END
exit 0
