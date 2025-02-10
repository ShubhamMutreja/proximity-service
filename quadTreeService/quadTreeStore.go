package quadtreeservice

import (
	"math"
	"proximityService/models"
)

type Node struct {
	//Rectangle origin
	startLatitude  float64
	startLongitude float64

	//height and width of a Rectangle
	width  float64
	height float64

	//Store list of Businesses
	listOfBusinesses []models.BusinessSearch

	//QuadTree's 4 Children or 4 equally divided rectangles
	topLeftChild     *Node
	topRightChild    *Node
	bottomLeftChild  *Node
	bottomRightChild *Node
}

type QuadTree struct {
	quadTreeNode              *Node
	maxSizeAllowedAtLeafNodes int
}

// Use haversine instead of euclidian distance for longitudes and latitudes
func haversine(sourceLatitudeDegrees, sourceLongitudeDegrees, destLatitudeDegrees, destLongitudeDegrees float64) float64 {
	//Radius of Earth
	const R = 6371

	sourceLatitudeRadians := sourceLatitudeDegrees * (math.Pi / 180)
	sourceLongitudeRadians := sourceLongitudeDegrees * (math.Pi / 180)
	destLatitudeRadians := destLatitudeDegrees * (math.Pi / 180)
	destLongitudeRadians := destLongitudeDegrees * (math.Pi / 180)

	dLat := destLatitudeRadians - sourceLatitudeRadians
	dLon := destLongitudeRadians - sourceLongitudeRadians

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(sourceLatitudeRadians)*math.Cos(destLatitudeRadians)*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}

func GetDistance(sourceLocation, destLocation models.Location) float64 {
	d := haversine(sourceLocation.Latitude, sourceLocation.Longitude, destLocation.Latitude, destLocation.Longitude)
	return d
}

func intersects(n *Node, loc models.Location, radius float64) bool {
	// Find closest point in the node to the search location
	closestLat := math.Max(n.startLatitude, math.Min(loc.Latitude, n.startLatitude+n.width))
	closestLon := math.Max(n.startLongitude, math.Min(loc.Longitude, n.startLongitude+n.height))

	// Compute distance between user location and the closest point in bounding box
	return haversine(loc.Latitude, loc.Longitude, closestLat, closestLon) <= radius
}

func NewQuadTree(latitude, longitude, width, height float64, maxSize int) *QuadTree {
	var node Node
	node.startLatitude = latitude
	node.startLongitude = longitude
	node.width = width
	node.height = height
	return &QuadTree{
		quadTreeNode:              &node,
		maxSizeAllowedAtLeafNodes: maxSize,
	}
}

func (quadTreeNode *Node) isLocationWithinBoundingBox(location models.Location) bool {
	if location.Latitude >= float64(quadTreeNode.startLatitude) &&
		location.Latitude < float64(quadTreeNode.startLatitude)+float64(quadTreeNode.width) &&
		location.Longitude >= float64(quadTreeNode.startLongitude) &&
		location.Longitude < float64(quadTreeNode.startLongitude)+float64(quadTreeNode.height) {
		return true
	}
	return false
}

func (quadTreeNode *Node) divideIntoFourBoxes() {
	//starting point of the parent bounding box
	startLatitude := quadTreeNode.startLatitude
	startLongitude := quadTreeNode.startLongitude

	//Dividing parent node into 4 children
	newHeight := quadTreeNode.height / 2
	newWidth := quadTreeNode.width / 2

	quadTreeNode.topLeftChild = &Node{
		startLatitude:  startLatitude,
		startLongitude: startLongitude,
		width:          newWidth,
		height:         newHeight,
	}
	quadTreeNode.topRightChild = &Node{
		startLatitude:  startLatitude + newWidth,
		startLongitude: startLongitude,
		width:          quadTreeNode.width - newWidth,
		height:         newHeight,
	}
	quadTreeNode.bottomLeftChild = &Node{
		startLatitude:  startLatitude,
		startLongitude: startLongitude + newHeight,
		width:          newWidth,
		height:         quadTreeNode.height - newHeight,
	}
	quadTreeNode.bottomRightChild = &Node{
		startLatitude:  startLatitude + newWidth,
		startLongitude: startLongitude + newHeight,
		width:          quadTreeNode.width - newWidth,
		height:         quadTreeNode.height - newHeight,
	}
}

func (quadTreeNode *Node) canNewNodeBeInserted(business models.BusinessSearch) bool {
	if quadTreeNode.isLocationWithinBoundingBox(business.Location) {
		return true
	}
	return false
}

func (quadTreeNode *Node) InsertNewNode(businessData models.BusinessSearch) {
	if !quadTreeNode.isLocationWithinBoundingBox(businessData.Location) {
		return
	}

	//Node is not a leaf. So we will check for which children we can add more nodes
	if !quadTreeNode.isThisNodeChild() {
		if quadTreeNode.topLeftChild.canNewNodeBeInserted(businessData) {
			quadTreeNode.topLeftChild.InsertNewNode(businessData)
		}
		if quadTreeNode.topRightChild.canNewNodeBeInserted(businessData) {
			quadTreeNode.topRightChild.InsertNewNode(businessData)
		}
		if quadTreeNode.bottomLeftChild.canNewNodeBeInserted(businessData) {
			quadTreeNode.bottomLeftChild.InsertNewNode(businessData)
		}
		if quadTreeNode.bottomRightChild.canNewNodeBeInserted(businessData) {
			quadTreeNode.bottomRightChild.InsertNewNode(businessData)
		}
		return
	}

	//node is a leaf.
	if len(quadTreeNode.listOfBusinesses) < 5 {
		quadTreeNode.listOfBusinesses = append(quadTreeNode.listOfBusinesses, businessData)
		return
	}
	// too many businesses on this leaf node, time to divide this node to 4 children
	quadTreeNode.divideIntoFourBoxes()
	quadTreeNode.listOfBusinesses = append(quadTreeNode.listOfBusinesses, businessData)
	for _, e := range quadTreeNode.listOfBusinesses {
		quadTreeNode.topLeftChild.InsertNewNode(e)
		quadTreeNode.topRightChild.InsertNewNode(e)
		quadTreeNode.bottomLeftChild.InsertNewNode(e)
		quadTreeNode.bottomRightChild.InsertNewNode(e)
	}
	//clear the list as only leaf nodes have this list
	quadTreeNode.listOfBusinesses = []models.BusinessSearch{}
}

func (quadTreeNode *Node) GetNearbyEntitiesFromQuadTree(userLoc models.Location, searchRadius float64) []models.BusinessSearch {
	var results []models.BusinessSearch

	if quadTreeNode == nil {
		return results
	}

	// Check if this node's bounding box intersects with the search radius
	if !intersects(quadTreeNode, userLoc, searchRadius) {
		return results
	}

	// Check businesses in this node
	for _, business := range quadTreeNode.listOfBusinesses {
		if haversine(userLoc.Latitude, userLoc.Longitude, business.Location.Latitude, business.Location.Longitude) <= searchRadius {
			results = append(results, business)
		}
	}

	// Recursively check child nodes
	results = append(results, quadTreeNode.topLeftChild.GetNearbyEntitiesFromQuadTree(userLoc, searchRadius)...)
	results = append(results, quadTreeNode.topRightChild.GetNearbyEntitiesFromQuadTree(userLoc, searchRadius)...)
	results = append(results, quadTreeNode.bottomLeftChild.GetNearbyEntitiesFromQuadTree(userLoc, searchRadius)...)
	results = append(results, quadTreeNode.bottomRightChild.GetNearbyEntitiesFromQuadTree(userLoc, searchRadius)...)

	return results
}

func (quadTreeNode *Node) isThisNodeChild() bool {
	return quadTreeNode.topLeftChild == nil && quadTreeNode.topRightChild == nil &&
		quadTreeNode.bottomLeftChild == nil && quadTreeNode.bottomRightChild == nil
}

func (qT *QuadTree) GetNearbyEntities(req models.NearbySearchRequest) []models.BusinessSearch {
	userLocation := req.UserLocation
	return qT.quadTreeNode.GetNearbyEntitiesFromQuadTree(userLocation, req.Radius)
}

func (qT *QuadTree) UpdateQuadTree(listOfBusinessesInDB []models.Business) {
	for _, businessDB := range listOfBusinessesInDB {
		var businessQT models.BusinessSearch
		businessQT.Location = businessDB.Location
		businessQT.BusinessID = businessDB.ID
		businessQT.Dist = GetDistance(businessQT.Location, models.Location{
			Longitude: 76.903872,
			Latitude:  28.842158,
		})
		qT.quadTreeNode.InsertNewNode(businessQT)
	}
}
