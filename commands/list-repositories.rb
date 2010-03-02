project_name = read_project_name
authorize(project_name, 'read')
repositories = []
Dir.foreach(find_project_dir(project_name)) do |repository_dir|
  next if repository_dir =~ /^\./
  git_url = "#{ENV['USER']}@#{Hubbard::HUB_HOST}:#{project_name}/#{repository_dir}"
  repositories << { :name => repository_dir.chomp('.git'), :url => git_url }
end

if @options[:format] == :yaml
  puts YAML::dump(repositories)
else
  repositories.each { |r| puts "#{r[:name]}\t#{r[:url]}" }
end
