projects = []
Dir.foreach(File.join(Hubbard::HUB_DATA, 'projects')) do |dir|
  next if dir == "." || dir == ".."
  if is_authorized(dir, 'read')
    projects << dir
  end
end

if OPTIONS[:format] == :yaml
  puts YAML::dump(projects)
else
  projects.each { |p| puts p }
end
