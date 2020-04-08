package com.ka.loaneligibility.backend.demo

import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.PostMapping
import org.springframework.web.bind.annotation.RequestBody
import org.springframework.web.bind.annotation.RequestMapping
import org.springframework.web.bind.annotation.RestController

@RestController
@RequestMapping("/demo")
class DemoCalculatorController(val calculator: DemoCalculator) {
    @GetMapping
    fun demo(): String = "Hello Calculator"

    @PostMapping("/calculator/plus")
    fun plus(@RequestBody body: plusInput): Int {
        return calculator.plus(body.a, body.b)
    }
}

data class plusInput(val a: Int = 0, val b: Int = 0)
