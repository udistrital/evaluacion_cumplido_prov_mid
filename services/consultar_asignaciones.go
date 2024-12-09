package services

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/evaluacion_cumplido_prov_mid/helpers"
	"github.com/udistrital/evaluacion_cumplido_prov_mid/models"
)

func ObtenerListaDeAsignaciones(documento string) (mapResponse map[string]interface{}, outputError error) {

	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()

	mapResponse = make(map[string]interface{})
	var listaAsignacionEvaluador []models.AsignacionEvaluador
	var listaAsignaciones []models.AsignacionEvaluacion
	listaAsignacionEvaluador, err := consultarAsignaciones(documento)

	if err != nil {
		return nil, fmt.Errorf("Error al consultar asignaciones")
	}

	if len(listaAsignacionEvaluador) > 0 && listaAsignacionEvaluador[0].EvaluacionId != nil {
		for _, asignacion := range listaAsignacionEvaluador {

			listaContratoGeneral, err := obtenerContratoGeneral(asignacion.EvaluacionId.ContratoSuscritoId, asignacion.EvaluacionId.VigenciaContrato)

			if err != nil {
				return nil, fmt.Errorf("Error al consultar detalles del contrato")

			}

			if len(listaContratoGeneral) > 0 {
				var contratoGeneral = listaContratoGeneral[0]

				respuesta, err := obtenerProveedor(contratoGeneral.Contratista, asignacion, listaContratoGeneral)
				listaAsignaciones = append(listaAsignaciones, respuesta)
				if err != nil {
					return nil, fmt.Errorf("Error al consultar detalles del contrato")

				}

			}

		}

	}

	dependencias, err := obtenerDependencias(documento)
	if err != nil {
		return nil, fmt.Errorf("Error al consultar asignaciones")
	}

	var listaContratosSupervisor []models.Contrato

	for _, dependencia := range dependencias {
		contrato_dependencia, _ := consultarContratosPorDependencia(dependencia.Codigo)
		listaContratosSupervisor = append(listaContratosSupervisor, contrato_dependencia...)
	}

	var contratosDepedencia []models.ContratoGeneral

	for _, contrato := range listaContratosSupervisor {

		numeroContrato, _ := strconv.Atoi(contrato.NumeroContrato)
		numeroVigencia, _ := strconv.Atoi(contrato.Vigencia)

		contratoGeneral, _ := obtenerContratoGeneral(numeroContrato, numeroVigencia)
		contratosDepedencia = append(contratosDepedencia, contratoGeneral...)
	}

	var listaSinAsignaciones []models.AsignacionEvaluacion
	for _, contrato := range contratosDepedencia {

		var listaProveedor []models.Proveedor

		if response, err := helpers.GetJsonWSO2Test(beego.AppConfig.String("UrlcrudAgora")+"/informacion_proveedor/?query=Id:"+strconv.Itoa(contrato.Contratista), &listaProveedor); err == nil && response == 200 {
			asignacionEvaluacion := models.AsignacionEvaluacion{
				AsignacionEvaluacionId: 0,
				NombreProveedor:        listaProveedor[0].NomProveedor,
				Dependencia:            contrato.DependenciaSolicitante,
				TipoContrato:           contrato.TipoContrato.TipoContrato,
				NumeroContrato:         contrato.ContratoSuscrito[0].NumeroContratoSuscrito,
				VigenciaContrato:       strconv.Itoa(contrato.VigenciaContrato),
				EvaluacionId:           0,
				Estado:                 false,
			}
			listaSinAsignaciones = append(listaSinAsignaciones, asignacionEvaluacion)

		} else {
			return nil, fmt.Errorf("Error al consultar asignaciones")

		}
	}

	if err != nil {
		return nil, fmt.Errorf("Error al consultar asignaciones")
	}

	mapResponse["Asignaciones"] = limpiarSinAsignaciones(listaAsignaciones, listaSinAsignaciones)
	mapResponse["SinAsignaciones"] = listaSinAsignaciones
	return mapResponse, nil
}

func consultarAsignaciones(documento string) (asignaciones []models.AsignacionEvaluador, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()
	var respuestaPeticion map[string]interface{}
	var listaAsignacionEvaluador []models.AsignacionEvaluador

	if response, err := helpers.GetJsonWSO2Test(beego.AppConfig.String("urlEvaluacionCumplidosCrud")+"/asignacion_evaluador?query=personaId:"+documento, &respuestaPeticion); err == nil && response == 200 {
		helpers.LimpiezaRespuestaRefactor(respuestaPeticion, &listaAsignacionEvaluador)
		if len(listaAsignacionEvaluador) > 0 && listaAsignacionEvaluador[0].EvaluacionId != nil {
			asignaciones = listaAsignacionEvaluador

		}
	} else {
		return asignaciones, fmt.Errorf("Error al consultar asignaciones")

	}
	return asignaciones, nil
}

func obtenerContratoGeneral(contratoSuscritoId int, vigenciaContrato int) (contratoGeneral []models.ContratoGeneral, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()

	if response, err := helpers.GetJsonWSO2Test(beego.AppConfig.String("UrlcrudAgora")+"/contrato_general/?query=ContratoSuscrito.NumeroContratoSuscrito:"+strconv.Itoa(contratoSuscritoId)+",VigenciaContrato:"+strconv.Itoa(vigenciaContrato), &contratoGeneral); err == nil && response == 200 {
	} else {
		return contratoGeneral, fmt.Errorf("Error al consultar asignaciones")

	}
	return contratoGeneral, nil
}

func obtenerContratoGeneralPorNumeroDecontrato(contratoSuscritoId int, vigenciaContrato int) (contratoGeneral []models.ContratoGeneral, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()

	if response, err := helpers.GetJsonWSO2Test(beego.AppConfig.String("UrlcrudAgora")+"/contrato_general/?query=ContratoSuscrito.Id:"+strconv.Itoa(contratoSuscritoId)+",VigenciaContrato:"+strconv.Itoa(vigenciaContrato), &contratoGeneral); err == nil && response == 200 {
	} else {
		return contratoGeneral, fmt.Errorf("Error al consultar asignaciones")

	}
	return contratoGeneral, nil
}

func obtenerProveedor(contratistaId int, asignacion models.AsignacionEvaluador, listaContratoGeneral []models.ContratoGeneral) (asisgnaciones models.AsignacionEvaluacion, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()

	var listaProveedor []models.Proveedor
	contratoGeneral := listaContratoGeneral[0]

	if response, err := helpers.GetJsonWSO2Test(beego.AppConfig.String("UrlcrudAgora")+"/informacion_proveedor/?query=Id:"+strconv.Itoa(contratistaId), &listaProveedor); err == nil && response == 200 {
		asignacionEvaluacion := models.AsignacionEvaluacion{
			AsignacionEvaluacionId: asignacion.Id,
			NombreProveedor:        listaProveedor[0].NomProveedor,
			Dependencia:            contratoGeneral.DependenciaSolicitante,
			TipoContrato:           contratoGeneral.TipoContrato.TipoContrato,
			NumeroContrato:         contratoGeneral.ContratoSuscrito[0].NumeroContratoSuscrito,
			VigenciaContrato:       strconv.Itoa(contratoGeneral.VigenciaContrato),
			EvaluacionId:           asignacion.EvaluacionId.Id,
			Estado:                 asignacion.Activo,
		}
		asisgnaciones = asignacionEvaluacion
	} else {
		return asisgnaciones, fmt.Errorf("Error al consultar asignaciones")

	}

	return asisgnaciones, nil
}

func obtenerDependencias(documento string) (dependencias []models.Dependencia, outputError error) {

	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()

	var respuesta models.DependenciasRespuesta

	fmt.Println(beego.AppConfig.String("UrlAdministrativaJBPM") + "/dependencias_supervisor/" + documento)
	if response, err := helpers.GetJsonWSO2Test(beego.AppConfig.String("UrlAdministrativaJBPM")+"/dependencias_supervisor/"+documento, &respuesta); err == nil && response == 200 {
		dependencias = respuesta.Dependencias.Dependencia
	} else {
		return dependencias, fmt.Errorf("Error al consultar asignaciones")
	}
	return dependencias, nil
}

func consultarContratosPorDependencia(dependencia string) (contratos []models.Contrato, outputError error) {

	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()

	var respuesta models.ContratosRespuesta

	fmt.Println((beego.AppConfig.String("UrlAdministrativaJBPM") + "/contratos_proveedor_dependencia/" + dependencia))
	if response, err := helpers.GetJsonWSO2Test(beego.AppConfig.String("UrlAdministrativaJBPM")+"/contratos_proveedor_dependencia/"+dependencia, &respuesta); err == nil && response == 200 {
		contratos = respuesta.Contratos.Contrato
	} else {
		return contratos, fmt.Errorf("Error al consultar depenendiencias")
	}
	return contratos, nil

}

func limpiarSinAsignaciones(Asignaciones, SinAsignaciones []models.AsignacionEvaluacion) []models.AsignacionEvaluacion {
	asignacionesMap := make(map[string]bool)
	for _, a := range Asignaciones {
		key := fmt.Sprintf("%d-%s", a.AsignacionEvaluacionId, a.VigenciaContrato)
		asignacionesMap[key] = true
	}

	var filtroSinAsignaciones []models.AsignacionEvaluacion
	for _, sa := range SinAsignaciones {
		key := fmt.Sprintf("%d-%s", sa.AsignacionEvaluacionId, sa.VigenciaContrato)
		if !asignacionesMap[key] {
			filtroSinAsignaciones = append(filtroSinAsignaciones, sa)
		}
	}

	return filtroSinAsignaciones
}
