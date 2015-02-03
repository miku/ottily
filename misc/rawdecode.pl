use strict;
use warnings;
use JSON qw( decode_json );

my $counter = 0; 
my $file = $ARGV[0];
open my $info, $file or die "Could not open $file: $!";

while(my $line = <$info>)  {
    my $decoded = decode_json($line);
    $counter++;
    if ($counter % 1000000 == 0) {
        print($counter . "\n");
    }
}

