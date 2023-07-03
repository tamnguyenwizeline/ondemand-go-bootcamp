package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/tamnguyenwizeline/ondemand-go-bootcamp/util"
)

func (server *Server) getPokemons(ctx *gin.Context) {

	pokemons, err := util.ReadCSVDataWithWorkers("./resources/pokemons.csv", 10, 2)
	if err != nil {

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, pokemons)
}

func (server *Server) loadPokemons(ctx *gin.Context) {

	err := util.LoadPokemonData("./resources/pokemons.csv")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, nil)
}
