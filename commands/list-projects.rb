projects_dir = File.join(Hubbard::HUB_DATA, 'projects')
projects = []
Dir.foreach(projects_dir) do |project|
  next if project == "." || project == ".."
  if is_authorized(project, 'read')
    vis = File.read(File.join(projects_dir, project, ".visibility")).strip 
    desc = File.read(File.join(projects_dir, project, ".description")).strip
    projects << { :name => project, :visibility => vis, :description => desc }
  end
end

if OPTIONS[:format] == :yaml
  puts YAML::dump(projects)
else
  projects.each do |p| 
    puts "#{p[:name].ljust(16)} #{p[:visibility].ljust(10)} #{p[:description]}"
  end
end
