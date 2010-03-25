project_name = read_project_name
authorize(project_name, 'admin')
project = Project.new(project_name)

permissions = project.lock { project.permissions }

if @options[:format] == :yaml
  puts YAML::dump(permissions)
elsif @options[:format] == :text
  permissions.each do |permission|
    puts "#{permission[:user]}=#{permission[:access]}"
  end
end
