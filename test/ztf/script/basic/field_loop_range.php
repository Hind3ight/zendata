#!/usr/bin/env php
<?php
/**
[case]
title=
cid=0
pid=0

[group]
 >>
 >>

[esac]
*/

include_once __DIR__ . DIRECTORY_SEPARATOR . '../lib/zd.php';

$zd = new zendata();
$output = $zd->create("", "test.yaml", 10, "", array("fields"=>"field_loop_range"));

$count = sprintf("%d", count($output));
print(">> $count\n");

print(">> $output[2]\n");