class Project
  attr_reader :project_name

  def initialize(project_name)
    @project_name = project_name
    @dir = find_project_dir project_name
  end

  def create
    if File.exist?(@dir)
      error 4, "Project already exists with that name"
    end

    unless Dir.mkdir(@dir)
      error 1, "Unable to create directory: #{dir}"
    end
  end

  def lock
    File.open(File.join(@dir, ".lock"), "w+") do |lock|
      lock.flock(File::LOCK_EX)
      begin
        yield
      ensure
        lock.flock(File::LOCK_UN)
      end
    end
  end

  def visibility=(visibility)
    File.open(File.join(@dir, ".visibility"), "w") { |f| f << "#{visibility}\n" }
  end
  
  def description=(description)
    File.open(File.join(@dir, ".description"), "w") { |f| f << "#{description}\n" }
  end

  def permissions
    if File.exists?(permissions_file)
      File.read(permissions_file).split("\n").map do |line|
        permission = line.strip.split('=')
        { :user => permission.first, :access => permission.last }
      end
    else
      []
    end  
  end

  def add_permission(username,action)
    new_permissions = permissions.select { |permission| permission[:user] != username }
    if username != 'admin'
      new_permissions << { :user => username, :access => action }
    end
    File.open(permissions_file, "w") do |file| 
      new_permissions.each { |permission| file << permission[:user] << '=' << permission[:access] << "\n" }
    end
  end

  def remove_permission(username)
    new_permissions = permissions.select { |permission| permission[:user] != username }
    File.open(permissions_file, "w") do |file| 
      new_permissions.each { |permission| file << permission[:user] << '=' << permission[:access] << "\n" }
    end
  end

  private

  def permissions_file
    File.join(@dir, ".permissions")
  end
end
