project_name = read_project_name
description = next_arg "Please specify a project description"

project = Project.new(project_name)
project.create
project.add_permission(@username,'admin')
project.visibility=@options[:private] ? "private" : "public"
project.description=description
