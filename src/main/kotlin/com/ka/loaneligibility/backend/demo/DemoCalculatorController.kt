package com.ka.loaneligibility.backend.demo

import org.slf4j.Logger
import org.slf4j.LoggerFactory
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.PostMapping
import org.springframework.web.bind.annotation.RequestBody
import org.springframework.web.bind.annotation.RequestMapping
import org.springframework.web.bind.annotation.RestController
import org.springframework.web.client.RestTemplate

@RestController
@RequestMapping("/demo")
class DemoCalculatorController(val calculator: DemoCalculator, val restTemplate: RestTemplate) {
    private val log: Logger = LoggerFactory.getLogger(DemoCalculatorController::class.java)

    @GetMapping
    fun demo(): String = "Hello Calculator"

    @PostMapping("/calculator/plus")
    fun plus(@RequestBody body: plusInput): Int {
        return calculator.plus(body.a, body.b)
    }

    @GetMapping("/get-vehicle")
    fun hello(): Vehicle? {
        val url = "http://loan-eligibility-svc:8080/demo/vehicle"
        var result: Vehicle? = restTemplate.getForObject<Vehicle>(url, Vehicle::class.java)
        result!!.miles = 2500
        return result
    }

    @GetMapping("/vehicle")
    fun vehicle(): Vehicle {
        return Vehicle()
    }
}

data class plusInput(val a: Int = 0, val b: Int = 0)
class Vehicle {
    var brand: String = "Subaru"
    var model: String = "WRX STI"
    var miles: Int = 1000
}
