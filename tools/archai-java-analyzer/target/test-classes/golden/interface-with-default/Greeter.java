package com.example.greet;

public interface Greeter {
    String greeting();

    default String hello(String name) {
        return greeting() + " " + name;
    }
}
