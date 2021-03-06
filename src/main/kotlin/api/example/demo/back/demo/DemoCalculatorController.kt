package api.example.demo.back.demo

import api.example.demo.back.models.CarModel
import api.example.demo.back.repositories.CarRepository
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
class DemoCalculatorController(val calculator: DemoCalculator, val restTemplate: RestTemplate, var carRepository: CarRepository) {

    private val log: Logger = LoggerFactory.getLogger(DemoCalculatorController::class.java)

    @GetMapping
    fun demo(): String = "Hello Calculator"

    @PostMapping("/calculator/plus")
    fun plus(@RequestBody body: plusInput): Int {
        return calculator.plus(body.a, body.b)
    }

    @GetMapping("/vehicle")
    fun vehicle(): Vehicle {
        return Vehicle()
    }

    @GetMapping("/findAll")
    fun findAll(): MutableIterable<CarModel> {
        return carRepository.findAll()
    }

    @PostMapping("/save")
    fun save(@RequestBody carModel: CarModel): String {
        carRepository.save(carModel)
        return "Done."
    }
}

data class plusInput(val a: Int = 0, val b: Int = 0)
class Vehicle {
    var brand: String = "Subaru"
    var model: String = "WRX STI"
    var miles: Int = 1000
}
