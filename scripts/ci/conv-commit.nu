#!/usr/bin/env nu

pwd;
ls -l ./scripts/agit;

use ../../scripts/agit
let is_conv: bool = agit is_conventional_commit

if $is_conv == true {
  print "Conventional Commit"
} else {
  print "Not a Conventional Commit"
}
