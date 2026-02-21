#!/usr/bin/env ruby
# frozen_string_literal: true

require 'json'
require 'open3'
require 'pathname'
require 'fileutils'

def run(cmd, cwd: nil, env: nil)
  if cwd.nil?
    out, status = Open3.capture2e(env || {}, *cmd)
  else
    out, status = Open3.capture2e(env || {}, *cmd, chdir: cwd.to_s)
  end
  return out if status.success?

  warn "error: command failed: #{cmd.join(' ')}"
  warn out
  exit 1
end

def iter_go_mod_dirs(root)
  Pathname.new(root).glob('**/go.mod').map(&:dirname).reject do |p|
    parts = p.each_filename.to_a
    parts.include?('vendor') || parts.include?('.git') || parts.include?('.idea') || parts.include?('.vscode')
  end
end

def each_json_object(str)
  depth = 0
  in_string = false
  escape = false
  start_idx = nil

  str.each_char.with_index do |ch, i|
    if in_string
      if escape
        escape = false
      elsif ch == '\\'
        escape = true
      elsif ch == '"'
        in_string = false
      end
      next
    end

    if ch == '"'
      in_string = true
      next
    end

    if ch == '{'
      start_idx = i if depth == 0
      depth += 1
      next
    end

    if ch == '}'
      depth -= 1
      if depth == 0 && !start_idx.nil?
        yield str[start_idx..i]
        start_idx = nil
      end
    end
  end
end

def iter_go_list_modules(module_dir, mod_mode)
  env = { 'GOFLAGS' => "-mod=#{mod_mode}" }
  out = run(['go', 'list', '-m', '-json', 'all'], cwd: module_dir, env: env)
  each_json_object(out) do |obj_str|
    yield JSON.parse(obj_str)
  end
end

def module_version_key(m)
  path = m['Path']
  version = m['Version']
  return nil if path.nil? || version.nil?

  "#{path}@#{version}"
end

def modcache_entries(modcache)
  Pathname.new(modcache).each_child.select(&:directory?).reject do |p|
    p.basename.to_s == 'cache' || !p.basename.to_s.include?('@')
  end
end

def delete_download_cache(modcache, module_path, version, dry_run)
  base = Pathname.new(modcache).join('cache', 'download', module_path, '@v')
  return 0 unless base.exist?

  removed = 0
  %w[.mod .zip .info .ziphash].each do |suffix|
    p = base.join("#{version}#{suffix}")
    next unless p.exist?

    removed += 1
    if dry_run
      puts "DRY-RUN delete #{p}"
    else
      p.delete
    end
  end

  lock = base.join("#{version}.lock")
  if lock.exist?
    removed += 1
    if dry_run
      puts "DRY-RUN delete #{lock}"
    else
      lock.delete
    end
  end

  removed
end

root = '.'
apply = false
modcache = nil
mod_mode = 'readonly'

ARGV.each_with_index do |arg, i|
  case arg
  when '--root'
    root = ARGV[i + 1]
  when '--apply'
    apply = true
  when '--modcache'
    modcache = ARGV[i + 1]
  when '--mod'
    mod_mode = ARGV[i + 1]
  end
end
mod_mode = 'readonly' if mod_mode.nil? || mod_mode.empty?

modcache ||= run(['go', 'env', 'GOMODCACHE']).strip
if modcache.empty?
  warn 'error: could not determine GOMODCACHE'
  exit 1
end

go_mod_dirs = iter_go_mod_dirs(root)
if go_mod_dirs.empty?
  warn 'error: no go.mod files found under root'
  exit 1
end

keep = {}

go_mod_dirs.each do |d|
  iter_go_list_modules(d, mod_mode) do |m|
    key = module_version_key(m)
    keep[key] = true if key
  end
end

dry_run = !apply
removed_dirs = 0
removed_downloads = 0

modcache_entries(modcache).each do |entry|
  name = entry.basename.to_s
  next if keep.key?(name)

  at = name.rindex('@')
  next if at.nil?

  module_path = name[0...at]
  version = name[(at + 1)..]

  removed_downloads += delete_download_cache(modcache, module_path, version, dry_run)

  if dry_run
    puts "DRY-RUN delete #{entry}"
  else
    FileUtils.rm_rf(entry.to_s)
  end
  removed_dirs += 1
end

mode = dry_run ? 'DRY-RUN' : 'APPLY'
puts "#{mode} done. module dirs removed: #{removed_dirs}, download cache files removed: #{removed_downloads}"
