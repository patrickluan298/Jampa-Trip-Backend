#!/bin/sh

if [ ! "$(docker network ls | grep host-jampa-trip-network)" ]; then
  echo "Network inexistente. Criando.."
  docker network create host-jampa-trip-network
else
  echo "host-jampa-trip-network network existente."
fi
