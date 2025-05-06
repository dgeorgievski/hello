# Check if last commit message complies with conventional commits specs
export def is_conventional [] {
  let subject = git show -s --format=%s 
  let body = git show -s --format=%b
  
}
