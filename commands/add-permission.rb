project_name = read_project_name
authorize(project_name, 'admin')
dir = find_project_dir(project_name)
username = ARGV.shift
action = ARGV.shift
validate_action_name(action)

File.open(File.join(dir, ".lock"), "w+") do |lock|
  lock.flock(File::LOCK_EX)
  begin
    filename = File.join(dir, ".permissions")
    permissions = File.read(filename).split("\n").map { |line| line.strip }.select { |line| line.split('=')[0] != username }
    permissions << "#{username}=#{action}"
    File.open(filename, "w") do |file| 
      permissions.each { |permission| file << permission << "\n" }
    end
  ensure
    lock.flock(File::LOCK_UN)
  end
end
