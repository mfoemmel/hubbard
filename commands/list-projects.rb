projects = []
Dir.foreach(Hubbard::PROJECTS_PATH) do |project|
  next if project == "." || project == ".."
  if is_authorized(project, 'read')
    project_path = File.join(Hubbard::PROJECTS_PATH, project)
    vis = File.read(File.join(project_path, ".visibility")).strip 
    desc = File.read(File.join(project_path, ".description")).strip
    projects << { :name => project, :visibility => vis, :description => desc }
  end
end

projects = projects.sort_by { |project| project[:name] }

if OPTIONS[:format] == :yaml
  puts YAML::dump(projects)
else
  projects.each do |p| 
    puts "#{p[:name].ljust(16)} #{p[:visibility].ljust(10)} #{p[:description]}"
  end
end
