project_name = read_project_name
authorize(project_name, 'read')
repositories = []
Dir.foreach(find_project_dir(project_name)) do |repository_name|
  next if repository_name =~ /^\./
  git_url = "#{ENV['USER']}@#{Hubbard::HUB_HOST}:#{project_name}/#{repository_name}.git"
  repositories << { :name => repository_name, :url => git_url }
end

if OPTIONS[:format] == :yaml
  puts YAML::dump(repositories)
else
  repositories.each { |r| puts "#{r[:name]}\t#{r[:url]}" }
end
