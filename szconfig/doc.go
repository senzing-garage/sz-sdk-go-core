/*
The szconfig package is used to modify the in-memory representation of a Senzing configuration.
It is a wrapper over Senzing's G2Config C binding.

To use szconfig,
the LD_LIBRARY_PATH environment variable must include a path to Senzing's libraries.
Example:

	export LD_LIBRARY_PATH=/opt/senzing/g2/lib
*/
package szconfig
