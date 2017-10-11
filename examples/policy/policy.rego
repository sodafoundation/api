package opa.policies

import data.pools
import data.docks
import data.profiles

find_dock_by_pool_id[poolid] = dock {
	pools[_] = pool
	pool.id = poolid
	docks[_] = dock
	pool.dockId = dock.id
}

find_supported_pools[profilename] = pool {
	pools[_] = pool
	pool.freeCapacity > desired_capacity[profilename]
	contains(pool.parameters.diskType, desired_diskType[profilename])
	pool.parameters.iops > desired_iops[profilename]
	pool.parameters.bandwidth > desired_bandwidth[profilename]
}

desired_capacity[profilename] = capacity {
	profiles[_] = profile
	profile.name = profilename
	capacity = profile.extra.capacity	
}

desired_diskType[profilename] = diskType {
	profiles[_] = profile
	profile.name = profilename
	not profile.extra.diskType
	diskType = ""
} {
	profiles[_] = profile
	profile.name = profilename
	profile.extra.diskType
	diskType = profile.extra.diskType
}

desired_iops[profilename] = iops {
	profiles[_] = profile
	profile.name = profilename
	not profile.extra.iops
	iops = 0
} {
	profiles[_] = profile
	profile.name = profilename
	profile.extra.iops
	iops = profile.extra.iops
}

desired_bandwidth[profilename] = bandwidth {
	profiles[_] = profile
	profile.name = profilename
	not profile.extra.bandwidth
	bandwidth = 0
} {
	profiles[_] = profile
	profile.name = profilename
	profile.extra.bandwidth
	bandwidth = profile.extra.bandwidth
}