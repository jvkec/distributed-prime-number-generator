// Worker implementation that handles prime number calculation tasks. Each worker
// connects to the coordinator, receives work assignments (number ranges), processes
// them using the specified algorithm, and returns results. Workers can operate
// independently across multiple machines.