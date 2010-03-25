project_name = read_project_name

project = Project.new(project_name)
project.create
project.add_permission(@username,'admin')
project.visibility=@options[:private] ? "private" : "public"
