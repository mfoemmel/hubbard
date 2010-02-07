project_name = read_project_name
authorize(project_name, 'admin')
dir = find_project_dir(project_name)
username = ARGV.shift

contents = ""
permissions_file = File.join(dir, ".permissions")
if File.exists?(permissions_file)
  File.open(permissions_file, "r+") do |f|
    f.flock(File::LOCK_EX)
    contents = f.read
    f.flock(File::LOCK_UN)
  end
end

if OPTIONS[:format] == :yaml
  permissions = []
  contents.split("\n").map do |l|
    l.strip!
    p = l.split('=')
    permissions << { :user => p.first, :access => p.last }
  end
  puts YAML::dump(permissions)
elsif OPTIONS[:format] == :text
  puts contents unless contents.empty?
end
