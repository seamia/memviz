# memviz 

this a (somewhat extensively) modified version of the original [memviz](https://github.com/bradleyjkemp/memviz)

## what changed
in no particular order:
* "better" shapes and orientation
* ability to specify colors (see config file)
* limited length of strings (config)
* limited number of elements in slices (config)
* ability to substiture int values with their names (config)
* ability to stop exposure at specified types/fields (config)
*  ability to exposure of "sensitive" data (config)
* config files to control behavior (global and local)

## usage
exactly the same as the original (just use diffrent `import`)

## config(uration) file(s)
see [this](https://raw.githubusercontent.com/seamia/memviz/master/memviz.options) example of a configuration file

the config file should be named `memviz.options`, and be located in one of the following two places:
* current directory of you app
* home directory of current user

## show me
go from this:
<p align="center">
  <img src="https://raw.githubusercontent.com/seamia/memviz/master/.media/before.svg">
</p>
to this:
<p align="center">
  <img src="https://raw.githubusercontent.com/seamia/memviz/master/.media/after.svg">
</p>
