project_name = read_project_name
dir = find_project_dir(project_name)
authorize(project_name, 'admin')
FileUtils.rm_rf(dir)
