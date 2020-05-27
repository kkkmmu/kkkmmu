#!/usr/bin/perl -w

use 5.010;
use warnings;
use diagnostics;

print "Hello, world!\n";

say "Hello, World!";

say "Hello" . "World!";

say "Hello" x 3;

say 5 x 4;

say 5 * 4;

say "5" + 4;

say "Z" . 5 * 7;

say "Z" . 5 x 7;

@lines = `perldoc -u -f atan2`;   #perldoc. lines is an array
foreach (@lines) {				  #Iterate over the array
	s/\w<([^>]+)>/\U$1/g;
	print;						  #Print each line
}

$int = 10;                        #all variable name start with $
$str = "12";
$intstr = $str x $int;
say $intstr;

print 1 lt 2;
print 1 gt 2;
print 1 eq 2;
print 1 eq 1;

if (1 eq 1) {
	print "true\n";
} else {
	print "false\n";
}

$line = <STDIN>;

if ($line eq "\n") {
	print "Empty line\n";
} else {
	print "$line \n";
}

$line = <STDIN>;

chomp($line);  #remove the trailer \n
if ($line eq "") {
	print "Empty line\n";
} else {
	print "$line \n";
}

chomp($line = <STDIN>);  #remove the trailer \n
if ($line eq "") {
	print "Empty line\n";
} else {
	print "$line \n";
}

$count = 1;

while ($count < 10) {
	$count += 1;
	print "$count\n";
}

$arr[0] = "test";
$arr[1] = 1;
$arr[2] = `pwd`;

print "$arr[0]\n";
print "$arr[1]\n";
print "$arr[2]\n";
print "$arr[-1]\n";
print "$arr[-2]\n";
print "$#arr\n";
print "$arr[$#arr]\n";

say (1..100) ;

say qw(rest, test, best);
say qw! rest, test, best !;
say qw# rest, test, best #;

($test, $best, $rest) = (1, "a", "bc");
say $test;
say $best;
say $rest;

@array = 5..9;
say "@array";
push(@array, 0);
push(@array, 100);
say "@array";
pop(@array);
say "@array";
push (@array, 98..100);
say "@array";

shift(@array);
say "@array";
unshift(@array, 1000);
say "@array";
unshift(@array, 1001..1005);
say "@array";
$array[1]=9999;
say "@array";

foreach $a (@array) {
	say $a.", ";
};

foreach (@array) {
	say $_.".";
};

$_ = "12kdafjaskdfjas";
say;
say reverse(@array);
say sort(@array);

$n = @array;
@list = @array;
say $n;
say @list;
