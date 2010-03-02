name = next_arg("Please specify the key name")
if name !~ Hubbard::KEY_NAME_REGEX
  error 1, "Not a valid key name (letters and numbers only)"
end

dirname = File.join(find_account_dir(@username), "keys")
FileUtils.mkdir_p(dirname)

filename = File.join(dirname, name)
if !File.exist?(filename)
  error 1, "Key not found"
end

unless FileUtils.rm(filename)
  error 1, "Unable to delete key"
end

sync_keys
