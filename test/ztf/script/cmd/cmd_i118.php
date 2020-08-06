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
 >>
 >>

[esac]
*/

include_once __DIR__ . DIRECTORY_SEPARATOR . '../lib/zd.php';

$zd = new zendata();

$zd->changeLang("en");

$output = $zd->cmd("-help");
print(">> $output[0]\n");
$output = $zd->cmd("-example");
print(">> $output[0]\n");

$zd->changeLang("zh");

$output = $zd->cmd("-help");
print(">> $output[0]\n");
$output = $zd->cmd("-example");
print(">> $output[0]\n");

$zd->changeLang("en");