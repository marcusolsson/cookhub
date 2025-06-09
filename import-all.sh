#!/bin/bash

config_files=(
    marcusolsson/recipes
	# cooklang/recipes
	# nicholaswilde/recipes
	# Net-Mist/remy
	# dubadub/cookbook
	# bubonicfred/cookbook
	# Diegothx/CookBook
	# javieruhk/RecipeManager
	# pubmania/a_diabetics_journal
	# 2xAA/cookbook
	# andoriyu/cooking
	# azlekov/recipebook
	# BraeTroutman/cookbook
	# brendanmckenzie/recipes
	# briansunter/site
	# crhuber/cooking
	# cvhooser/recipes
	# demosjarco/recipes
	# dzenzes/recipes
	# firefly2442/my-recipes-cooklang
	# ggalmazor/recipes
	# Gunde/recipes
	# iandennismiller/recipes
	# ignacio-gn/cook
	# isaacvando/recipes
	# jnobles/RecipeBook
	# JorenC/helanxiaochu
	# justintout/recipes
	# leMaik/rezepte
	# LiHRaM/cooklang-recipes
	# LudeeD/comidadaboa
	# luizribeiro/menu
	# ngalaiko/galaiko.rocks
	# openpaul/cookbook
	# shen-sat/cooklang
	# simcard0000/recipes
	# surzycki/recipes
	# tdstein/recipes
	# tntraina/recipe_archive
	# TyHil/recipes
	# XpiritBV/innovation-dinner-2022
)

for config in "${config_files[@]}"; do
    echo "Importing ${config}..."
    curl -X POST localhost:8080/api/github.com/${config}
done



