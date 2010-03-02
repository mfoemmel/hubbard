project_name = read_project_name
authorize(project_name, 'admin')
dir = find_project_dir(project_name)
username = ARGV.shift
action = ARGV.shift
validate_action_name(action)

project = Project.new(project_name)
project.lock { project.add_permission(username,action) }
