package com.ka.loaneligibility.backend

import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.runApplication
import org.springframework.cloud.client.discovery.EnableDiscoveryClient


@SpringBootApplication
@EnableDiscoveryClient
class LoanEligibilityServiceApplication

fun main(args: Array<String>) {
    runApplication<LoanEligibilityServiceApplication>(*args)
}
