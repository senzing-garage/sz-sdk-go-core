/*
The szconfigmanager package is used to modify Senzing configuration JSON documents in the Senzing datastore.
It is a wrapper over Senzing's G2Configmgr C binding.

To use szconfigmanager,
the LD_LIBRARY_PATH environment variable must include a path to Senzing's libraries.
Example:

	export LD_LIBRARY_PATH=/opt/senzing/g2/lib
*/
package szconfigmanager
