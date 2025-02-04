package rest_handlers

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/litmuschaos/litmus/litmus-portal/graphql-server/pkg/cluster"
	"github.com/litmuschaos/litmus/litmus-portal/graphql-server/pkg/database/mongodb"
	dbSchemaCluster "github.com/litmuschaos/litmus/litmus-portal/graphql-server/pkg/database/mongodb/cluster"
	dbOperationsWorkflow "github.com/litmuschaos/litmus/litmus-portal/graphql-server/pkg/database/mongodb/workflow"
	"github.com/litmuschaos/litmus/litmus-portal/graphql-server/pkg/k8s"
	"github.com/litmuschaos/litmus/litmus-portal/graphql-server/utils"
	log "github.com/sirupsen/logrus"
)

// FileHandler dynamically generates the manifest file and sends it as a response
func FileHandler(mongodbOperator mongodb.MongoOperator) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := strings.TrimSuffix(c.Param("key"), ".yaml")
		kubeCluster, err := k8s.NewKubeCluster()
		if err != nil {
			log.Fatalf("error while getting kube config: %v", err)
		}
		response, statusCode, err := cluster.NewService(
			dbSchemaCluster.NewClusterOperator(mongodbOperator),
			dbOperationsWorkflow.NewChaosWorkflowOperator(mongodbOperator),
			kubeCluster,
		).GetManifest(token)
		if err != nil {
			log.WithError(err).Error("error while generating manifest file")
			utils.WriteHeaders(&c.Writer, statusCode)
			c.Writer.Write([]byte(err.Error()))
			return
		}
		utils.WriteHeaders(&c.Writer, statusCode)
		c.Writer.Write(response)
	}
}
