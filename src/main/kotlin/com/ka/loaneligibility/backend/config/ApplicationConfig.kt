package com.ka.loaneligibility.backend.config

import org.springframework.cloud.client.loadbalancer.LoadBalanced
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.web.client.RestTemplate

@Configuration
class ApplicationConfig {
    @LoadBalanced
    @Bean
    fun restTemplate(): RestTemplate {
        return RestTemplate()
    }
}
