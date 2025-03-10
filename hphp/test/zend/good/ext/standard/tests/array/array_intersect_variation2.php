<?hh
/* Prototype  : array array_intersect(array $arr1, array $arr2 [, array $...])
 * Description: Returns the entries of arr1 that have values which are present in all the other arguments
 * Source code: ext/standard/array.c
*/

/*
* Testing array_intersect() function by passing values to $arr2 argument other than arrays
* and see that function emits proper warning messages wherever expected.
* The $arr1 argument is a fixed array.
*/

// get a class
class classA
{
  public function __toString() :mixed{
    return "Class A object";
  }
}
<<__EntryPoint>> function main(): void {
echo "*** Testing array_intersect() : Passing non-array values to \$arr2 argument ***\n";

// array to be passsed to $arr1 as default argument
$arr1 = vec[1, 2];

// arrays to be passed to optional argument
$arr3 = dict[0 => 1, 1 => 2, "one" => 1, "two" => 2];


// heredoc string
$heredoc = <<<EOT
hello world
EOT;

// get a resource variable
$fp = fopen(__FILE__, "r");

// unexpected values to be passed to $arr2 argument
$arrays = vec[

       // int data
/*1*/  0,
       1,
       12345,
       -2345,

       // float data
/*5*/  10.5,
       -10.5,
       12.3456789000e10,
       12.3456789000E-10,
       .5,

       // null data
/*10*/ NULL,
       null,

       // boolean data
/*12*/ true,
       false,
       TRUE,
       FALSE,

       // empty data
/*16*/ "",
       '',

       // string data
/*18*/ "string",
       'string',
       $heredoc,

       // object data
/*21*/ new classA(),



       // resource variable
/*22*/ $fp
];

// loop through each sub-array within $arrrays to check the behavior of array_intersect()
$iterator = 1;
foreach($arrays as $unexpected_value) {
  echo "\n-- Iterator $iterator --";

  // Calling array_intersect() with default arguments
  var_dump( array_intersect($arr1,$unexpected_value) );

  // Calling array_intersect() with more arguments
  var_dump( array_intersect($arr1, $unexpected_value, $arr3) );

  $iterator++;
}

// close the file resource used
fclose($fp);

echo "Done";
}
