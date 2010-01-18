$stderr.puts <<-END
Usage: hub <command>

Projects:

list-projects
create-project <project>
delete-project <project>

Repositories:

create-repository <project> <repository>
delete-repository <project> <repository>

Permissions:

add-permission <project> <username> read|write|admin
remove-permission <project> <username>

END
exit 0
