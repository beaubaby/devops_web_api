package api.example.demo.back.demo

import org.springframework.stereotype.Service

@Service
class DemoCalculator {
    fun plus(a: Int, b: Int): Int = a + b
}
