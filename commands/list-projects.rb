Dir.foreach(File.join(Hubbard::HUB_DATA, 'projects')) do |dir|
  next if dir == "." || dir == ".."
  if is_authorized(dir, 'read')
    puts dir
  end
end
