require 'fileutils'

project_name = read_project_name
authorize(project_name, 'admin')

description = rest_args "Please specify the description"
desc = description.join(' ')

dir = find_project_dir(project_name)
if !File.exist?(dir)
  error 4, "Project not found"
end

File.open(File.join(dir, ".description"), "w") { |f| f << desc << "\n" }

