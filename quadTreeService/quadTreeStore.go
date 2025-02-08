package quadtreeservice

import (
	"math"
	"proximityService/models"
)

type Node struct {
	//Box starting points
	startLatitude  int
	startLongitude int

	//height and width of a Rectangle
	width  int
	height int

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

func NewQuadTree(latitude, longitude, width, height, maxSize int) *QuadTree {
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

func (qT *QuadTree) isLocationWithinBoundingBox(location models.Location, quadTreeNode *Node) bool {
	if location.Latitude >= float64(quadTreeNode.startLatitude) &&
		location.Latitude < float64(quadTreeNode.startLatitude)+float64(quadTreeNode.width) &&
		location.Longitude <= float64(quadTreeNode.startLongitude) &&
		location.Longitude > float64(quadTreeNode.startLongitude)-float64(quadTreeNode.height) {
		return true
	}
	return false
}

func (qT *QuadTree) divideIntoFourBoxes(quadTreeNode *Node) {
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
		startLongitude: startLongitude - newHeight,
		width:          newWidth,
		height:         quadTreeNode.height - newHeight,
	}
	quadTreeNode.bottomRightChild = &Node{
		startLatitude:  startLatitude + newWidth,
		startLongitude: startLongitude - newHeight,
		width:          quadTreeNode.width - newWidth,
		height:         quadTreeNode.height - newHeight,
	}
}

func (qT *QuadTree) canNewNodeBeInserted(business models.BusinessSearch, quadTreeNodeChild *Node) bool {
	if quadTreeNodeChild.topLeftChild != nil && qT.isLocationWithinBoundingBox(business.Location, quadTreeNodeChild.topLeftChild) {
		return true
	}
	return false
}

func (qT *QuadTree) InsertNewNode(businessData models.BusinessSearch, quadTreeNode *Node) {
	if !qT.isLocationWithinBoundingBox(businessData.Location, quadTreeNode) {
		return
	}

	//Node is not a leaf. So we will check for which children we can add more nodes
	if !qT.isThisNodeChild(quadTreeNode) {
		if qT.canNewNodeBeInserted(businessData, quadTreeNode.topLeftChild) {
			qT.InsertNewNode(businessData, quadTreeNode.topLeftChild)
		}
		if qT.canNewNodeBeInserted(businessData, quadTreeNode.topRightChild) {
			qT.InsertNewNode(businessData, quadTreeNode.topRightChild)
		}
		if qT.canNewNodeBeInserted(businessData, quadTreeNode.bottomLeftChild) {
			qT.InsertNewNode(businessData, quadTreeNode.bottomLeftChild)
		}
		if qT.canNewNodeBeInserted(businessData, quadTreeNode.bottomRightChild) {
			qT.InsertNewNode(businessData, quadTreeNode.bottomRightChild)
		}
		return
	}

	//node is a leaf.
	if len(quadTreeNode.listOfBusinesses) < qT.maxSizeAllowedAtLeafNodes {
		quadTreeNode.listOfBusinesses = append(quadTreeNode.listOfBusinesses, businessData)
		return
	}
	// too many businesses on this leaf node, time to divide this node to 4 children
	qT.divideIntoFourBoxes(quadTreeNode)
	quadTreeNode.listOfBusinesses = append(quadTreeNode.listOfBusinesses, businessData)
	for _, e := range quadTreeNode.listOfBusinesses {
		qT.InsertNewNode(e, quadTreeNode.topLeftChild)
		qT.InsertNewNode(e, quadTreeNode.topRightChild)
		qT.InsertNewNode(e, quadTreeNode.bottomLeftChild)
		qT.InsertNewNode(e, quadTreeNode.bottomRightChild)
	}
	//clear the list as only leaf nodes have this list
	quadTreeNode.listOfBusinesses = []models.BusinessSearch{}
}

func (qT *QuadTree) UpdateQuadTree(listOfBusinessesInDB []models.Business) {
	for _, businessDB := range listOfBusinessesInDB {
		var businessQT models.BusinessSearch
		businessQT.Location = businessDB.Location
		businessQT.BusinessID = businessDB.ID
		qT.InsertNewNode(businessQT, qT.quadTreeNode)
	}
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

func getDistance(userLocation, businessLocation models.Location) float64 {
	return haversine(userLocation.Latitude, userLocation.Longitude, businessLocation.Latitude, businessLocation.Longitude)
}

func (qT *QuadTree) isLocationInsideRadius(node *Node, userLocation models.Location, radius float64) bool {
	//For topLeftChild
	x1 := float64(node.startLatitude)
	y1 := float64(node.startLongitude)
	dist := getDistance(models.Location{Longitude: x1, Latitude: y1}, userLocation)
	if dist <= radius {
		return true
	}

	//For topRightChild
	x1 = float64(node.startLatitude + node.width)
	y1 = float64(node.startLongitude)
	dist = getDistance(models.Location{Longitude: x1, Latitude: y1}, userLocation)
	if dist <= radius {
		return true
	}

	//For bottomLeftChild
	x1 = float64(node.startLatitude)
	y1 = float64(node.startLongitude - node.height)
	dist = getDistance(models.Location{Longitude: x1, Latitude: y1}, userLocation)
	if dist <= radius {
		return true
	}
	//For bottomRightChild
	x1 = float64(node.startLatitude + node.width)
	y1 = float64(node.startLongitude - node.height)
	dist = getDistance(models.Location{Longitude: x1, Latitude: y1}, userLocation)
	if dist <= radius {
		return true
	}
	return false
}

func (qT *QuadTree) isThisNodeChild(node *Node) bool {
	return (node.topLeftChild == nil && node.topRightChild == nil && node.bottomLeftChild == nil && node.bottomRightChild == nil)
}

func (qT *QuadTree) isUserLocationInsideThisNodeRadius(quadTreeNodeChild *Node, userLocation models.Location, radius float64) bool {
	if quadTreeNodeChild != nil &&
		(qT.isLocationInsideRadius(quadTreeNodeChild, userLocation, radius) || qT.isLocationWithinBoundingBox(userLocation, quadTreeNodeChild)) {
		return true
	}
	return false
}

func (qT *QuadTree) GetNearbyEntitiesFromQuadTree(userLocation models.Location, quadTreeNode *Node, radius float64) []models.BusinessSearch {
	var response []models.BusinessSearch
	if qT.isThisNodeChild(quadTreeNode) &&
		(qT.isLocationInsideRadius(quadTreeNode, userLocation, radius) ||
			qT.isLocationWithinBoundingBox(userLocation, quadTreeNode)) {
		listOfBusinessesAtLeafNode := quadTreeNode.listOfBusinesses
		for _, businessLocationData := range listOfBusinessesAtLeafNode {
			distFromUser := getDistance(businessLocationData.Location, userLocation)
			if distFromUser <= radius {
				businessLocationData.Dist = distFromUser
				response = append(response, businessLocationData)
			}
		}
	}

	//search in top left child
	if qT.isUserLocationInsideThisNodeRadius(quadTreeNode.topLeftChild, userLocation, radius) {
		tmpResponse := qT.GetNearbyEntitiesFromQuadTree(userLocation, quadTreeNode.topLeftChild, radius)
		response = append(response, tmpResponse...)
	}

	//search in top right child
	if qT.isUserLocationInsideThisNodeRadius(quadTreeNode.topRightChild, userLocation, radius) {
		tmpResponse := qT.GetNearbyEntitiesFromQuadTree(userLocation, quadTreeNode.topRightChild, radius)
		response = append(response, tmpResponse...)
	}

	//search in bottom Left child
	if qT.isUserLocationInsideThisNodeRadius(quadTreeNode.bottomLeftChild, userLocation, radius) {
		tmpResponse := qT.GetNearbyEntitiesFromQuadTree(userLocation, quadTreeNode.bottomLeftChild, radius)
		response = append(response, tmpResponse...)
	}

	//search in bottom right child
	if qT.isUserLocationInsideThisNodeRadius(quadTreeNode.bottomRightChild, userLocation, radius) {
		tmpResponse := qT.GetNearbyEntitiesFromQuadTree(userLocation, quadTreeNode.bottomRightChild, radius)
		response = append(response, tmpResponse...)
	}
	return response
}

func (qT *QuadTree) GetNearbyEntities(req models.NearbySearchRequest) []models.BusinessSearch {
	entitySearch := req.UserLocation
	return qT.GetNearbyEntitiesFromQuadTree(entitySearch, qT.quadTreeNode, req.Radius)
}
