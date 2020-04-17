package com.ka.loaneligibility.backend.demo

import org.slf4j.Logger
import org.slf4j.LoggerFactory
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.PostMapping
import org.springframework.web.bind.annotation.RequestBody
import org.springframework.web.bind.annotation.RequestMapping
import org.springframework.web.bind.annotation.RestController
import org.springframework.web.client.RestTemplate

@RestController
@RequestMapping("/demo")
class DemoCalculatorController(val calculator: DemoCalculator) {
    private val log: Logger = LoggerFactory.getLogger(DemoCalculatorController::class.java)
    @Autowired
    private val restTemplate: RestTemplate? = null
    @GetMapping
    fun demo(): String = "Hello Calculator"

    @PostMapping("/calculator/plus")
    fun plus(@RequestBody body: plusInput): Int {
        return calculator.plus(body.a, body.b)
    }

    @GetMapping("/hello-world")
    fun hello() {
        val url = "http://hello-service/"
        val result: String? = restTemplate!!.getForObject<String>(url, String::class.java)
        log.debug("Result: {}", result)
    }
}

data class plusInput(val a: Int = 0, val b: Int = 0)
