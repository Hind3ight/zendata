#!/usr/bin/env php
<?php
/**
[case]
title=指定步长
cid=1340
pid=7

[group]
 显示10行生成的数据 >>
 验证第3行数据     >>

[esac]
*/

include_once __DIR__ . DIRECTORY_SEPARATOR . '../lib/zd.php';

$zd = new zendata();
$output = $zd->create("", "advanced.yaml", 10, "", array("fields"=>"field_step_negative"));

$count = sprintf("%d", count($output));
print(">> $count\n");

print(">> $output[2]\n");