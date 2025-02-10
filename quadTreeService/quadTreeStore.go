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
	listOfBusinesses map[string]models.BusinessSearch

	//QuadTree's 4 Children or 4 equally divided rectangles
	listOfChildren []*Node
}

type QuadTree struct {
	quadTreeNode              *Node
	maxSizeAllowedAtLeafNodes int
}

// Use haversine instead of euclidian distance for longitudes and latitudes
func haversine(sourceLatitudeDegrees, sourceLongitudeDegrees, destLatitudeDegrees, destLongitudeDegrees float64) float64 {
	//Radius of Earth
	const earthRadiusKm = 6371

	sourceLatitudeRadians := sourceLatitudeDegrees * (math.Pi / 180)
	sourceLongitudeRadians := sourceLongitudeDegrees * (math.Pi / 180)
	destLatitudeRadians := destLatitudeDegrees * (math.Pi / 180)
	destLongitudeRadians := destLongitudeDegrees * (math.Pi / 180)

	deltaLat := destLatitudeRadians - sourceLatitudeRadians
	deltaLon := destLongitudeRadians - sourceLongitudeRadians

	havFormula := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(sourceLatitudeRadians)*math.Cos(destLatitudeRadians)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)

	centralAngle := 2 * math.Atan2(math.Sqrt(havFormula), math.Sqrt(1-havFormula))

	return earthRadiusKm * centralAngle
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
	listOfChildrenNodes := make([]*Node, 4)

	listOfChildrenNodes[0] = &Node{
		startLatitude:  startLatitude,
		startLongitude: startLongitude,
		width:          newWidth,
		height:         newHeight,
	}
	listOfChildrenNodes[1] = &Node{
		startLatitude:  startLatitude + newWidth,
		startLongitude: startLongitude,
		width:          quadTreeNode.width - newWidth,
		height:         newHeight,
	}
	listOfChildrenNodes[2] = &Node{
		startLatitude:  startLatitude,
		startLongitude: startLongitude + newHeight,
		width:          newWidth,
		height:         quadTreeNode.height - newHeight,
	}
	listOfChildrenNodes[3] = &Node{
		startLatitude:  startLatitude + newWidth,
		startLongitude: startLongitude + newHeight,
		width:          quadTreeNode.width - newWidth,
		height:         quadTreeNode.height - newHeight,
	}
	quadTreeNode.listOfChildren = listOfChildrenNodes
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
		for _, childNode := range quadTreeNode.listOfChildren {
			if childNode.canNewNodeBeInserted(businessData) {
				childNode.InsertNewNode(businessData)
			}
		}
		return
	}
	//takes care of panics in case map is nil
	if quadTreeNode.listOfBusinesses == nil {
		emptyMap := make(map[string]models.BusinessSearch)
		quadTreeNode.listOfBusinesses = emptyMap
	}
	//node is a leaf.
	if len(quadTreeNode.listOfBusinesses) < 5 {
		if _, ok := quadTreeNode.listOfBusinesses[businessData.BusinessID]; !ok {
			quadTreeNode.listOfBusinesses[businessData.BusinessID] = businessData
		}
		return
	}
	// too many businesses on this leaf node, time to divide this node to 4 children
	quadTreeNode.divideIntoFourBoxes()
	if _, ok := quadTreeNode.listOfBusinesses[businessData.BusinessID]; !ok {
		quadTreeNode.listOfBusinesses[businessData.BusinessID] = businessData
	}
	for _, e := range quadTreeNode.listOfBusinesses {
		for _, childNode := range quadTreeNode.listOfChildren {
			childNode.InsertNewNode(e)
		}
	}
	//clear the map as only leaf nodes have this filled
	emptyMap := make(map[string]models.BusinessSearch)
	quadTreeNode.listOfBusinesses = emptyMap
	return
}

func (quadTreeNode *Node) DeleteNode(businessData models.BusinessSearch) bool {
	if quadTreeNode.isThisNodeChild() {
		for _, business := range quadTreeNode.listOfBusinesses {
			if business.BusinessID == businessData.BusinessID {
				delete(quadTreeNode.listOfBusinesses, business.BusinessID)
				return true
			}
		}
		return false
	}
	var res bool
	for _, childNode := range quadTreeNode.listOfChildren {
		if childNode != nil && childNode.isLocationWithinBoundingBox(businessData.Location) {
			res = res || childNode.DeleteNode(businessData)
		}
	}
	return res

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
		distFromUser := haversine(userLoc.Latitude, userLoc.Longitude, business.Location.Latitude, business.Location.Longitude)
		if distFromUser <= searchRadius {
			business.Dist = distFromUser
			results = append(results, business)
		}
	}

	// Recursively check child nodes
	for _, childNode := range quadTreeNode.listOfChildren {
		results = append(results, childNode.GetNearbyEntitiesFromQuadTree(userLoc, searchRadius)...)
	}

	return results
}

func (quadTreeNode *Node) isThisNodeChild() bool {
	return quadTreeNode.listOfChildren == nil || len(quadTreeNode.listOfChildren) == 0
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
		qT.quadTreeNode.InsertNewNode(businessQT)
	}
}

func (qT *QuadTree) DeleteFromQuadTree(listOfBusinessesInDB []models.Business) {
	for _, businessDB := range listOfBusinessesInDB {
		var businessQT models.BusinessSearch
		businessQT.Location = businessDB.Location
		businessQT.BusinessID = businessDB.ID
		qT.quadTreeNode.DeleteNode(businessQT)
	}
}
