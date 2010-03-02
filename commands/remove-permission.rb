project_name = read_project_name
authorize(project_name, 'admin')
dir = find_project_dir(project_name)
username = ARGV.shift

project = Project.new(project_name)
project.lock { project.remove_permission(username) }
