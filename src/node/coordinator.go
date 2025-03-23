// Coordinator service that manages the distribution of work among connected worker
// nodes. It divides the requested prime number range into chunks, assigns these
// chunks to available workers, tracks progress, and collects results. Implements
// basic load balancing and worker failure handling.