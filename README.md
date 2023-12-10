# Movie Recommendation System in Go

This recommendation system is based on an Edx Data Science Course. The Movielens dataset of 10M movies and ratings was used. Based on the edx course the dataset was split to the edx dataset and the test dataset. Further, only the edx dataset was used which was further split 4:1 to the train set and the test set. Parameters of the models were computed using the train set only and then, the models were evaluated using the test set. RMSE and MAE were calculated.

Data containing directors and 2 main actors were web-scraped from csfd.cz.

## Modelv0 (0 variables)

For comparison purposes, the first developed model is the most simplest one:

```math
$$y_{ij} = a$$
```
where $a$ is the average of all ratings.

Summary statistics of this model based on test set are:

RMSE = 1.06
MAE = 0.86

## Modelv2 (2 variables)

```math
$$y_{ij} = a + m_{i} + u_{j} $$
```
where $m_{ij}$ is an effect of movie $i$ calculated as an average of $y_{ij} - a$ for every $j$ and 
$m_{ij}$ is an effect of user $j$ calculated as an average of $y_{ij} - a - m_{ij}$ for every $i$. Model can be trained with coefficient lambda which is added to the denominator of the fraction when calculating the average.

Summary statistics of this model based on test set given lambda = 0 are: 

RMSE = 0.87
MAE = 0.67

## Modelv4 (4 variables)

Model 
```math
$$y_{ij} = a + d_{i} + (a_{1i} + a_{2i})/2 + m_{i} + u_{j} $$
```
were $d_{i}$ is the director's effect, $a_{1i}$ and $a_{2i}$ are the effects of two main actors.

Summary statistics of this model based on test set given lambda = 0 are: 

RMSE = 0.87
MAE = 0.67


