package org.example;

import org.junit.jupiter.api.Test;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;

import static org.junit.jupiter.api.Assertions.*;

class MobilityConnectorTest {
    @Test
    public void listAllStationTest() throws IOException, InterruptedException {
        RouterResponse routerResponse = MobilityConnector.listAllStations();
        assertEquals(200,routerResponse.getStatus());
        assertEquals(readAllFileContent("mobility-flat-node.json"), routerResponse.getBody());
    }

    public static String readAllFileContent(String filePath) throws IOException {
        Path path = Paths.get(filePath);
        byte[] bytes = Files.readAllBytes(path);
        return new String(bytes);
    }
}